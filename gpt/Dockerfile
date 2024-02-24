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

# Copy dependencies
COPY requirements.txt /app

# Install dependencies
RUN pip install -r requirements.txt

# Copy the Python script into the container
COPY main.py /app

# Run the Python script when the container starts
# Use default parameters for a larger model
ENTRYPOINT ["python", "main.py" \
	,"--batch-size", "64" \
	,"--block-size", "256" \
	,"--blocks", "6" \
	,"--learning-rate", "3e-4" \
	,"--head-size", "24" \
	,"--heads", "6" \
	,"--steps", "5000" \
	,"--tokeniser-file", "/out/tokeniser.json" \
	,"--model-file", "/out/model.onnx" \
	]
