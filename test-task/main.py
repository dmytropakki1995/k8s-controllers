from flask import Flask, request, send_from_directory
import os
import mimetypes
import sqlite3
import time

app = Flask(__name__)
UPLOAD_FOLDER = os.getenv("UPLOAD_FOLDER", "./uploads") # "/app/uploads") # the default path os '/app/uploads' instead of /'uploads'
os.makedirs(UPLOAD_FOLDER, exist_ok=True)

db_path = "./app/data/uploads.db"
os.makedirs(os.path.dirname(db_path), exist_ok=True)

# Initialize the SQLite database
def init_db():
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    cursor.execute('''CREATE TABLE uploads (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        filename TEXT NOT NULL,
        upload_time TEXT NOT NULL,
        mime_type TEXT NOT NULL
    )''')
    conn.commit()
    conn.close()

@app.route('/health')
def health():
    return 'OK', 200

@app.route('/uploads/<filename>')
def get_file(filename):
    file_path = os.path.join(UPLOAD_FOLDER, filename)

    # Check if the file exists
    if not os.path.isfile(file_path):
        return {"Error": f"File {filename} does not exist in {UPLOAD_FOLDER}"}, 404
    
    # Send the file
    return send_from_directory(UPLOAD_FOLDER, filename)

@app.route('/upload', methods=['POST'])
def upload_file():
    if 'file' not in request.files:
        return 'No file part', 400
    file = request.files['file']
    if file.filename == '':
        return 'No selected file', 400

    # Determine file path and save it
    file_path = os.path.join(UPLOAD_FOLDER, file.filename)
    file.save(file_path)

    # Set file permissions
    os.chmod(file_path, 0o660)

    # Get MIME type
    mime_type = mimetypes.guess_type(file_path)[0] or 'unknown'
    print(f"MIME type: {mime_type}")

    
    # This method should not use init_db() every time when the file is uploaded. It fails if this method is called more than once.
    # init_db() is better to use in a separate method or even a container like a migration container.
    # Store metadata in the database
    init_db()

    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    cursor.execute('''INSERT INTO uploads (filename, upload_time, mime_type) VALUES (?, ?, ?)''',
                   (file.filename, time.strftime('%Y-%m-%d %H:%M:%S'), mime_type))
    conn.commit()
    conn.close()

    return f'File {file.filename} uploaded successfully', 200

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080)