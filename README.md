# GopherNotes

A basic Go+React application for helping with reviewing old long notes

I built this since I was revisiting old notes that were way too long to re-read in its entirety. I was also playing around with a locally hosted LLM on my machine at the time so I decided to make this application. A little bonus is that my notes never leave my local network so my docs stay private.

I hosted the 9B model on my PC using `llama.cpp` then set up a SSH Tunnel (Local port forwarding) to connect to the model securely on my laptop.

### Tech Stack
Backend: Go (Gin)\
Frontend: React (Vite)\
AI: Llama.cpp + Qwopus 3.5 9B (GGUF)

Link to model: https://huggingface.co/Jackrong/Qwopus3.5-9B-v3-GGUF