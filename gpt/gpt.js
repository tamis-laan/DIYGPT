// Sleep in milliseconds
function sleep(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

// Load the tokenizer
async function tokeniser(path) {
	// Fetch the tokenizer
	const response = await fetch(path);
	if (!response.ok) {
		throw new Error('Network response was not ok');
	}

	// Get the json data
	const data = await response.json();

	// Construct encoder
	const encode = s => Array.from(s, c => data.stoi[c]);

	// Construct decoder
	const decode = l => l.map(i => data.itos[i]).join('');

	// Return encoder decoder
	return { encode, decode };
}

// Add value to the end while the first value drops off
function enque(array, value) {
	for (let i = 0; i < array.length - 1; i++) {
		array[i] = array[i + 1]
	}
	array[array.length - 1] = value
}

// Generate new token
async function generate(session, input) {
	// Run the model
	const inputs = { input: input }; // Replace 'inputName' with your model's input name
	// Get the results
	const result = await session.run(inputs);
	// Access the output (adjust 'outputName' as needed)
	const output = result.output; // Replace 'outputName' with your model's output name
	// Return int
	return output.data[0]
}

// Attach model for generation
async function attach(button, output, block_size = 32, max_tokens = 1000, tokeniser_file = "tokeniser.json", model_file = "model.onnx") {

	// Load tokeniser
	const { _, decode } = await tokeniser(tokeniser_file)

	// Load the model
	const session = await ort.InferenceSession.create(model_file);

	// Event handler for run model button
	button.addEventListener('click', async function() {

		// Make output empty
		output.innerText = ""

		// Prepare the input tensor
		input = new ort.Tensor('int32', new Int32Array(block_size).fill(0), [block_size]);

		// Disable the button while generating text
		button.disabled = true

		// Generate text
		for (let i = 1; i <= max_tokens; i++) {
			// Generate a new token
			token_new = await generate(session, input)
			// Decode token
			token_decoded = decode([token_new])
			// Append new data to the end
			enque(input.data, token_new)
			// Append
			output.innerText += decode([token_new])
			await sleep(1)
		}

		// Enable the button
		button.disabled = false
	})

}
