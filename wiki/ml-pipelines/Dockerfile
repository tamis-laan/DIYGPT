FROM python:3.9-slim

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

# USER python

WORKDIR /app

COPY jobs /app/jobs

RUN ls /app
RUN ls /app/jobs
