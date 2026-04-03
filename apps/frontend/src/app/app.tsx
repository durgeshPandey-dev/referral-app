import { useState, useRef } from "react";
import { postV1Upload } from "../api/v1";

interface HistoryItem {
  id: string;
  name: string;
  status: string;
  time: string;
}

interface CardProps {
  title: string;
  value: string;
}

function App() {
  const [file, setFile] = useState<File | null>(null);
  const [message, setMessage] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [history, setHistory] = useState<HistoryItem[]>([]);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFile = (f: File | undefined) => {
    if (!f) return;

    // Validate file type
    const validTypes = [
      "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
      "application/vnd.ms-excel",
    ];

    if (!validTypes.includes(f.type) && !f.name.endsWith(".xlsx")) {
      setMessage("⚠️ Please select a valid Excel file (.xlsx)");
      return;
    }

    setFile(f);
    setMessage("");
  };

  const handleUpload = async () => {
    if (!file) {
      setMessage("⚠️ Please select a file");
      return;
    }

    setLoading(true);

    const formData = new FormData();
    formData.append("file", file);

    postV1Upload(
      { file },
      { baseURL: import.meta.env.VITE_API_URL || "http://localhost:3000/api" },
    )
      .then((data) => {
        setMessage("✅ " + (data?.message || "Processing started"));

        const newHistoryItem: HistoryItem = {
          id: `${file.name}-${Date.now()}`,
          name: file.name,
          status: "Processing",
          time: new Date().toLocaleTimeString(),
        };

        setHistory((prev) => [newHistoryItem, ...prev]);
      })
      .catch((e) => {
        console.error("Upload error:", e);
        setMessage("❌ Upload failed");
      })
      .finally(() => {
        setLoading(false);
      });
  };

  return (
    <div style={styles.container}>
      {/* Sidebar */}
      <div style={styles.sidebar}>
        <h2>
          <span role="img" aria-label="lightning">
            ⚡
          </span>
          Referral
        </h2>
        <p>Dashboard</p>
        <p>Uploads</p>
        <p>Analytics</p>
      </div>

      {/* Main Content */}
      <div style={styles.main}>
        <h1>
          <span role="img" aria-label="dashboard">
            📊{" "}
          </span>{" "}
          Dashboard
        </h1>

        {/* Stats */}
        <div style={styles.stats}>
          <Card title="Emails Sent" value="120" />
          <Card title="In Progress" value="5" />
          <Card title="Failures" value="2" />
        </div>

        {/* Upload Section */}
        <div style={styles.card}>
          <h3>
            <span role="img" aria-label="upload">
              📤
            </span>{" "}
            Upload Excel
          </h3>

          <div
            style={styles.dropZone}
            onDragOver={(e) => e.preventDefault()}
            onDrop={(e) => {
              e.preventDefault();
              const droppedFile = e.dataTransfer.files?.[0];
              handleFile(droppedFile);
            }}
            onClick={() => fileInputRef.current?.click()}
          >
            {file ? (
              <p>
                <span role="img" aria-label="file">
                  📄
                </span>
                {file.name}
              </p>
            ) : (
              <p>Drag & Drop Excel or Click</p>
            )}

            <input
              ref={fileInputRef}
              type="file"
              accept=".xlsx"
              onChange={(e) => {
                const selectedFile = e.target.files?.[0];
                handleFile(selectedFile);
              }}
              style={styles.hiddenInput}
            />
          </div>

          <button onClick={handleUpload} style={styles.button}>
            {loading ? "Uploading..." : "Upload"}
          </button>

          {message && <p>{message}</p>}
        </div>

        {/* Activity */}
        <div style={styles.card}>
          <h3>
            <span role="img" aria-label="activity">
              📜
            </span>{" "}
            Recent Activity
          </h3>

          {history.length === 0 && <p>No uploads yet</p>}

          {history.map((item) => (
            <div key={item.id} style={styles.activityItem}>
              <span>{item.name}</span>
              <span>{item.status}</span>
              <span>{item.time}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

const Card = ({ title, value }: CardProps) => (
  <div style={styles.statCard}>
    <h4>{title}</h4>
    <p style={styles.statValue}>{value}</p>
  </div>
);

const styles: { [key: string]: React.CSSProperties } = {
  container: {
    display: "flex",
    height: "100vh",
    fontFamily: "Arial",
  },
  sidebar: {
    width: "200px",
    background: "#111",
    color: "#fff",
    padding: "20px",
  },
  main: {
    flex: 1,
    padding: "30px",
    background: "#f5f7fa",
  },
  stats: {
    display: "flex",
    gap: "20px",
    marginBottom: "20px",
  },
  statCard: {
    background: "#fff",
    padding: "20px",
    borderRadius: "10px",
    flex: 1,
    boxShadow: "0 5px 15px rgba(0,0,0,0.1)",
  },
  statValue: {
    fontSize: "24px",
    fontWeight: "bold",
  },
  card: {
    background: "#fff",
    padding: "20px",
    borderRadius: "10px",
    marginBottom: "20px",
    boxShadow: "0 5px 15px rgba(0,0,0,0.1)",
  },
  dropZone: {
    border: "2px dashed #aaa",
    padding: "30px",
    textAlign: "center",
    marginBottom: "10px",
    cursor: "pointer",
    borderRadius: "10px",
  },
  hiddenInput: {
    display: "none",
  },
  button: {
    padding: "10px 20px",
    background: "#4CAF50",
    color: "#fff",
    border: "none",
    borderRadius: "5px",
    cursor: "pointer",
  },
  activityItem: {
    display: "flex",
    justifyContent: "space-between",
    padding: "10px 0",
    borderBottom: "1px solid #eee",
  },
};

export default App;
