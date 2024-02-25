# Use the official Python image as base
FROM python:latest

# Install dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
	cmake \
	&& rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Create output directory
RUN mkdir /out

# Mount as volume to get the tokeniser and model output
# EXAMPLE: docker run -v "$(pwd)/out:/out" gpt --steps 2500
VOLUME ["/out"]

# Copy dependencies
COPY requirements.txt /app

# Install dependencies
RUN pip install -r requirements.txt

# Copy the Python script into the container
COPY main.py /app

# Run the Python script when the container starts
# Use default parameters for a larger model
# Mount /tmp in RAM memory for speed gains
ENTRYPOINT ["sh", "-c", "mount -t tmpfs none /tmp && python main.py \
	--batch-size 64 \
	--block-size 128 \
	--blocks 6 \
	--learning-rate 3e-4 \
	--head-size 64 \
	--heads 6 \
	--steps 10000 \
	--tokeniser-file /out/tokeniser.json \
	--model-file /out/model.onnx" \
	]

