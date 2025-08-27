FROM python:3.12-slim

WORKDIR /app

ENV UPLOAD_FOLDER="/app/uploads"

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

EXPOSE 8080

# Run the Flask app
CMD ["python", "main.py"]