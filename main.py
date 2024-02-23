import os
import requests
import torch
import json

# Set the torch seed
torch.manual_seed(1337)

# Download dataset
def dataset():
    filename = "input.txt"
    url = "https://raw.githubusercontent.com/karpathy/char-rnn/master/data/tinyshakespeare/input.txt"
    # Download dataset
    if not os.path.exists(filename):
        response = requests.get(url)
        if response.status_code == 200:
            with open(filename, 'wb') as file:
                file.write(response.content)
        else:
            raise ValueError("Failed to download tinyshakespeare dataset")
    # Load dataset from disk
    with open(filename, 'r', encoding='utf-8') as file:
        content = file.read()
    return content

# Generate tokeniser from dataset
def tokenise(dataset):
    # Extract unique characters used
    vocab = sorted(list(set(dataset)))
    # Construct string to int mapping
    stoi = { ch:i for i,ch in enumerate(vocab) }
    # Construct int to string mapping
    itos = { i:ch for i,ch in enumerate(vocab) }
    # Construct tokeniser encoder
    encode = lambda s: [stoi[c] for c in s]
    # Construct tokeniser decoder
    decode = lambda l: ''.join([itos[i] for i in l])
    # Define tokeniser for export
    tokeniser = {
        "vocab":vocab,
        "stoi": stoi,
        "itos": itos,
    }
    # Return encoder decoder
    return encode, decode, vocab, tokeniser

# Sample a training batch from source dataset
def batch(source, block_size=8, batch_size=4):
    # Generate random sample positions
    ix = torch.randint(len(source)-block_size, (batch_size,))
    # Sample chunks
    x  = torch.stack([source[i:i+block_size] for i in ix])
    # Sample chunks with offset of 1
    y  = torch.stack([source[i+1:i+block_size+1] for i in ix])
    return x,y

# Self Attention Head
class SelfAttentionHead(torch.nn.Module):
    def __init__(self, T, C, H=16, encoder=True):
        super().__init__()
        self.T = T # Block size
        self.C = C # Channels
        self.H = H # Head size
        self.encoder    = encoder
        self.key   = torch.nn.Linear(C, H, bias=False)
        self.query = torch.nn.Linear(C, H, bias=False)
        self.value = torch.nn.Linear(C, H, bias=False)
        self.register_buffer('t', torch.tril(torch.ones(T,T)))
    def forward(self, x):
        # Compute the key
        k = self.key(x) # B,T,H
        # Compute the query
        q = self.query(x) # B,T,H
        # Compute the weights
        # (B,T,H) @ (B,H,T) --> (B,T,T)
        w = q @ k.transpose(-2,-1) * self.H**-0.5 
        # Mask the weights
        if self.encoder:
            w = w.masked_fill(self.t[:self.T,:self.T] == 0, float('-inf'))
        # Apply the softwax 
        w = torch.nn.functional.softmax(w, dim=-1)
        # Return result
        return w @ self.value(x)

class MultiSelfAttentionHead(torch.nn.Module):
    def __init__(self, T, C, H=16, heads=10, encoder=True):
        super().__init__()
        self.heads = torch.nn.ModuleList([SelfAttentionHead(T,C,H,encoder) for _ in range(heads)])

    def forward(self, x):
        return torch.cat([h(x) for h in self.heads], dim=-1)

class FeedForward(torch.nn.Module):
    def __init__(self, C):
        super().__init__()
        self.net = torch.nn.Sequential(
            torch.nn.Linear(C,4*C), 
            torch.nn.ReLU(),
            torch.nn.Linear(4*C,C), 
        )
    def forward(self,x):
        return self.net(x)

class Block(torch.nn.Module):
    def __init__(self, T, C, H=16, heads = 10, encoder = True):
        super().__init__()
        self.ln1  = torch.nn.LayerNorm(C)
        self.head = MultiSelfAttentionHead(T,C,H, heads, encoder)
        self.ln2  = torch.nn.LayerNorm(C)
        self.ff   = FeedForward(H*heads)
    def forward(self,x):
        x = x + self.head(self.ln1(x))
        x = x + self.ff(self.ln2(x))
        return x

# Bigram language model
class GPT(torch.nn.Module):

    def __init__(self, vocab_size, block_size, blocks=3, head_size = 16, heads=2):
        # Init super
        super().__init__()
        # Store parameters
        self.vocab_size = vocab_size
        self.block_size = block_size
        self.head_size  = head_size
        self.heads      = heads
        # Create bigram
        self.word_embedding = torch.nn.Embedding(vocab_size, head_size*heads)
        # Encode token position
        self.position_embedding = torch.nn.Embedding(block_size, head_size*heads)
        # Create self attention heads
        # self.blocks = Block(block_size, head_size, heads)
        self.blocks = torch.nn.Sequential(
            *[Block(block_size, head_size*heads, head_size, heads) for _ in range(blocks)]
        )
        # Layer norm for last layer
        self.ln = torch.nn.LayerNorm(head_size*heads)
        # Create linear layer for head
        self.last = torch.nn.Linear(head_size*heads, vocab_size)

    def forward(self, idx, targets=None):
        # Get dims
        B,T = idx.shape
        # Embed tokens
        token_embedding = self.word_embedding(idx) # (B,T,C)
        # Token position embedding
        position_embedding = self.position_embedding(torch.arange(T)) # (T,C)
        # Cat emebddings
        x = token_embedding + position_embedding # (B,T,C)
        # Apply self attention
        x = self.blocks(x)
        # Get logits
        logits = self.last(self.ln(x)) # (B,T, vocab_size)
        # Return
        if targets is None:
            return logits, None
        # Get dims
        B, T, C = logits.shape
        # Reshape logits
        rlogits = logits.view(B*T, C)
        # Reshape targets
        targets = targets.view(B*T)
        # Compute the loss
        loss = torch.nn.functional.cross_entropy(rlogits,targets.long())
        # Return 
        return logits, loss

    def generate(self, x):
        # Crop input
        x_cropped = x[:, -self.block_size:]
        logits, _ = self(x_cropped)
        logits = logits[:,-1,:]
        probs  = torch.nn.functional.softmax(logits, dim=1)
        x_next = torch.multinomial(probs, num_samples=1)
        return x_next.int()

# Train Bigram Language
def train_gpt(model, train, block_size=8, batch_size=32, steps=1000, lr=1e-3):
    # Set model to training mode
    model.train()
    # Smoothing decay
    decay = .95
    # Smoothing loss
    smooth_loss = None
    # Create optimizer
    optimizer = torch.optim.AdamW(model.parameters(), lr=lr)
    # Start training
    for step in range(steps):
        xb, yb = batch(train, block_size, batch_size)
        _, loss = model(xb,yb)
        if smooth_loss is None:
            smooth_loss = loss
        else:
            smooth_loss = decay*smooth_loss + (1-decay)*loss.item()
        print(f"\rstep: {step}/{steps} loss: {smooth_loss:.2f} {loss:.2f}", end="", flush=True)
        optimizer.zero_grad(set_to_none=True)
        loss.backward()
        optimizer.step()

# Main
if __name__ == "__main__":
    # Load dataset
    text = dataset()
    # Create tokeniser
    encode,decode,vocab, tokeniser = tokenise(text)
    # Save tokeniser to disk
    with open("tokeniser.json", "w") as json_file:
        json_file.write(json.dumps(tokeniser))
    # Encode text
    data = torch.tensor(encode(text), dtype=torch.int)
    # Split dataset
    n = int(0.9*len(data))
    # Training dataset
    train = data[:n]
    # Validation dataset
    val = data[n:]
    # Create bigram model
    model = GPT(len(vocab), 32, 3, 32, 8)
    # Train model
    train_gpt(model, train, batch_size=64, block_size=32, steps=2000, lr=1e-3)
    # Start with new line character
    idx = torch.zeros((1,32), dtype=torch.int)
    # Generate words
    model.eval()
    for _ in range(1000):
        idx  = torch.cat((idx, model.generate(idx)), dim=1)
        print(decode([idx[0][-1].tolist()]),end="")


    # Export model
    class OnnxModel(torch.nn.Module):
        def __init__(self, model):
            super().__init__()
            self.model = model
        def forward(self,x):
            # Add batch dimension on input
            x = x.unsqueeze(0)
            # Run generate
            y = self.model.generate(x)
            # Remove batch dimension
            y = y.squeeze(0)
            return y

    # Wrap the model for onnx export
    onnx_model = OnnxModel(model)

    # Export model in onnx format
    torch.onnx.export(
        onnx_model,
        torch.zeros((32), dtype=torch.int),
        "model.onnx", 
        input_names=['input'],
        output_names=['output']
    )

