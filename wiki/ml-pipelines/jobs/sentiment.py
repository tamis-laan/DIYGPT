import uvicorn
from transformers import pipeline
from fastapi import FastAPI
from pydantic import BaseModel

# Load the sentiment analysis pipeline
sentiment_pipeline = pipeline("sentiment-analysis")

# Create API
app = FastAPI()

class Body(BaseModel):
    sentence: str

# Create predict route
@app.post("/run")
async def run(body: Body):
    # Perform sentiment analysis
    result = sentiment_pipeline(body.sentence)
    # Return result
    return result

# Start server
if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
