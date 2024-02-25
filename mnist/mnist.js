// Function to display results
function displayResults(probabilities, container) {
	// Clear the container
	container.innerHTML = '';

	// Create a new canvas
	const canvas = document.createElement('canvas');
	const ctx = canvas.getContext('2d');
	canvas.width = 256;
	canvas.height = 128;

	// Draw probabilities
	const barWidth = canvas.width / probabilities.length;
	const maxProb = Math.max(...probabilities);
	const maxIndex = probabilities.indexOf(maxProb);

	for (let i = 0; i < probabilities.length; i++) {
		const barHeight = (probabilities[i] / maxProb) * (canvas.height - 40);
		const x = i * barWidth;
		const y = canvas.height - barHeight - 15;

		// Color the bar with max probability red
		ctx.fillStyle = i === maxIndex ? 'red' : 'blue';
		ctx.fillRect(x, y, barWidth, barHeight);

		// Draw text indicating the probability
		ctx.fillStyle = 'black';
		ctx.fillText(probabilities[i].toFixed(2), x + barWidth / 2 - 10, y - 10);

		// Draw text indicating the digit below the bars
		ctx.fillText(i, x + barWidth / 2 - 5, canvas.height + 15);
	}

	// Draw digit numbers
	ctx.fillStyle = 'black';
	for (let i = 0; i < 10; i++) {
		const x = i * barWidth + barWidth / 2 - 5;
		const y = canvas.height;
		ctx.fillText(i, x, y);
	}

	// Append visualisation
	container.appendChild(canvas);
}


// Attach
async function attach(canvas, classify, clear, results, model_file) {

	// Disable buttons while loading
	classify.disabled = true
	clear.disabled = true

	// Display loading text
	results.innerText = 'Loading model...'

	// Load the model
	const session = await ort.InferenceSession.create(model_file);

	// Clear loading text
	results.innerText = ''

	// Enable buttons again
	classify.disabled = false
	clear.disabled = false

	// Get canvas context
	const ctx = canvas.getContext('2d');

	// Is drawing global
	let isDrawing = false;

	// Functions for drawing
	function startDrawing(e) {
		isDrawing = true;
		draw(e);
	}

	function draw(e) {
		if (!isDrawing) return;
		ctx.lineWidth = 25;
		ctx.lineCap = 'round';
		ctx.strokeStyle = 'white';
		ctx.lineTo(e.offsetX, e.offsetY);
		ctx.stroke();
		ctx.beginPath();
		ctx.moveTo(e.offsetX, e.offsetY);
	}

	function stopDrawing() {
		isDrawing = false;
		ctx.beginPath();
	}

	// Function to clear canvas
	function clearCanvas() {
		ctx.clearRect(0, 0, canvas.width, canvas.height);
		results.innerText = ''
	}

	// Event listeners for drawing
	canvas.addEventListener('mousedown', startDrawing);
	canvas.addEventListener('mousemove', draw);
	canvas.addEventListener('mouseup', stopDrawing);
	canvas.addEventListener('mouseout', stopDrawing);

	// Event listener for classification
	clear.addEventListener('click', clearCanvas)

	// Event listener for clearing canvas
	classify.addEventListener('click', async function() {

		// Disable the button while classifying
		classify.disabled = true

		// Show loading on results
		results.innerText = 'processing ...'

		// Create a dummy canvas just for scaling
		const scaledCanvas = document.createElement('canvas');

		// Get the context
		const scaledCtx = scaledCanvas.getContext('2d');

		// Scale the image
		scaledCtx.drawImage(canvas, 0, 0, 28, 28);

		// Get the image data
		data = scaledCtx.getImageData(0, 0, 28, 28).data

		// Convert to grayscale and normalize
		inputarray = []
		for (let i = 0; i < data.length; i += 4) {
			// Turn into grayscale
			grayscaled = (data[i] + data[i + 1] + data[i + 2]) / 3.0
			// Scale between 0 and 1
			downscaled = grayscaled / 255.0
			// Normalize subtracting mean and variance
			normalized = (downscaled - 0.1307) / 0.3081
			// Push onto array
			inputarray.push(normalized);
		}

		// Create input tensor
		tensor = new ort.Tensor('float32', new Float32Array(inputarray), [1, 1, 28, 28]);

		// Run the model
		const inputs = { input: tensor };
		// Get the results
		const result = await session.run(inputs);
		// Access the output (adjust 'outputName' as needed)
		const output = result.output;
		// Create probabilities
		probs = output.data.map(x => Math.exp(x))

		// Disable the button while generating text
		classify.disabled = false

		// Draw results
		displayResults(probs, results)
	})

}
