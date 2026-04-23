import { useState, useEffect, useRef } from 'react'
import './App.css'
import ReactMarkdown from 'react-markdown'

const API = import.meta.env.VITE_API_URL || "http://localhost:8080"

function App() {
  const [status, setStatus] = useState("checking...")
  const [noteStatus, setNoteStatus] = useState(null)
  const [noteError, setNoteError] = useState(false)
  const [messages, setMessages] = useState([])
  const [input, setInput] = useState("")
  const [loading, setLoading] = useState(false)
  const bottomRef = useRef(null)

  useEffect(() => {
    // test backend connection
    fetch(`${API}/ping`)
      .then(res => res.json())
      .then(data => setStatus(data.message))
      .catch(() => setStatus("unreachable"))
  }, [])

  // Auto-scroll to bottom when messages change
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" })
  }, [messages, loading])

  async function handleUpload(e) {
    const file = e.target.files[0]
    if (!file) return

    const form = new FormData()
    form.append("note", file)

    try {
      const res = await fetch(`${API}/upload`, { method: "POST", body: form })
      const data = await res.json()

      if (!res.ok) {
        setNoteStatus(`Error: ${data.error}`)
        setNoteError(true)
        return
      }

      setNoteStatus(`${data.filename} · ${data.chars.toLocaleString()} chars`)
      setNoteError(false)
    } catch {
      setNoteStatus("Could not reach backend")
      setNoteError(true)
    }
  }

  async function sendMessage() {
    if (!input.trim() || loading) return

    const userMessage = input.trim()
    setInput("")
    setLoading(true)
    setMessages(prev => [...prev, { role: "user", text: userMessage }])

    try {
      const res = await fetch(`${API}/chat`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ message: userMessage }),
      })
      const data = await res.json()

      if (!res.ok) {
        setMessages(prev => [...prev, { role: "bot", text: `Error: ${data.error}` }])
        return
      }

      setMessages(prev => [...prev, { role: "bot", text: data.reply }])
    } catch {
      setMessages(prev => [...prev, { role: "bot", text: "Could not reach backend" }])
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      {/* Header */}
      <div className="header">
        <div className="header-left">
          <div className="header-logo">🦫</div>
          <span className="header-title">GopherNotes</span>
        </div>
      </div>

      {/* Upload bar */}
      <div className="upload-bar">
        <label className="upload-label">
          + Upload note (.md or .pdf)
          <input type="file" accept=".md,.pdf" onChange={handleUpload} />
        </label>
        {noteStatus && (
          <span className={`note-badge ${noteError ? "error" : ""}`}>
            {noteStatus}
          </span>
        )}
      </div>

      {/* Messages */}
      <div className="messages">
        {messages.length === 0 && !loading && (
          <div className="empty-state">Upload a note and start chatting</div>
        )}

      {messages.map((m, i) => (
        <div key={i} className={`message-row ${m.role}`}>
          {m.role === "bot" && <div className="bot-avatar">🦫</div>}
          <div className={`bubble ${m.role}`}>
            {m.role === "bot"
              ? <ReactMarkdown>{m.text}</ReactMarkdown>
              : m.text
            }
          </div>
        </div>
      ))}

        {loading && (
          <div className="message-row bot">
            <div className="bot-avatar">🦫</div>
            <div className="bubble bot">
              <div className="typing">
                <span /><span /><span />
              </div>
            </div>
          </div>
        )}

        <div ref={bottomRef} />
      </div>

      {/* Input */}
      <div className="input-bar">
        <input
          value={input}
          onChange={e => setInput(e.target.value)}
          onKeyDown={e => e.key === "Enter" && sendMessage()}
          placeholder="Ask something about your note..."
          disabled={loading}
        />
        <button className="send-btn" onClick={sendMessage} disabled={loading}>
          {loading ? "..." : "Send"}
        </button>
      </div>
    </>
  )
}

export default App