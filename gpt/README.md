# Do it yourself GPT

This is a mini GPT model written from scratch using pytorch.

Note that you can break out of the training loop using ctl-C and the script will save the model to disk from the last training step.

Create a virtual environment
``` bash
python -m venv .venv
```

Install required packages
``` bash
pip install -r requirements.txt
```

Activate the environment
``` bash
source .venv/bin/activate
```

Run the script
``` bash
python main.py
```

You can also use docker
``` bash
docker build . --tag gpt:latest
```

``` bash
docker run -v "$(pwd)/out:/out" gpt:latest --steps 2500
```

Serve the demo website
``` bash
python serve.py
```

