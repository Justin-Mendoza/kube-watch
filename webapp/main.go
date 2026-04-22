package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
)

const page = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>kube-watch</title>
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body {
    background: #0f0f1a;
    color: #e0e0e0;
    font-family: 'Courier New', monospace;
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    overflow: hidden;
  }
  .container {
    text-align: center;
    z-index: 1;
    position: relative;
  }
  h1 {
    font-size: 3rem;
    background: linear-gradient(135deg, #326ce5, #00d2ff);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    margin-bottom: 0.5rem;
  }
  .status {
    display: inline-flex;
    align-items: center;
    gap: 10px;
    background: rgba(50, 108, 229, 0.1);
    border: 1px solid rgba(50, 108, 229, 0.3);
    border-radius: 999px;
    padding: 8px 20px;
    margin-top: 1rem;
    font-size: 0.9rem;
  }
  .dot {
    width: 10px;
    height: 10px;
    background: #00e676;
    border-radius: 50%;
    animation: pulse 1.5s ease-in-out infinite;
  }
  @keyframes pulse {
    0%, 100% { opacity: 1; box-shadow: 0 0 0 0 rgba(0, 230, 118, 0.6); }
    50% { opacity: 0.6; box-shadow: 0 0 0 12px rgba(0, 230, 118, 0); }
  }

  /* floating kubernetes wheels */
  .wheel {
    position: fixed;
    opacity: 0.06;
    animation: float linear infinite;
  }
  @keyframes float {
    0%   { transform: translateY(110vh) rotate(0deg); }
    100% { transform: translateY(-20vh) rotate(360deg); }
  }

  /* pod grid */
  .pods {
    display: flex;
    gap: 12px;
    justify-content: center;
    margin-top: 2rem;
  }
  .pod {
    width: 48px;
    height: 48px;
    border-radius: 10px;
    background: linear-gradient(135deg, #326ce5, #1b3a80);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.7rem;
    color: rgba(255,255,255,0.8);
    animation: popIn 0.4s ease-out both;
    position: relative;
    overflow: hidden;
  }
  .pod::after {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(135deg, transparent 40%, rgba(255,255,255,0.1));
    border-radius: 10px;
  }
  .pod:nth-child(1) { animation-delay: 0.1s; }
  .pod:nth-child(2) { animation-delay: 0.2s; }
  .pod:nth-child(3) { animation-delay: 0.3s; }
  .pod:nth-child(4) { animation-delay: 0.4s; }
  @keyframes popIn {
    0%   { transform: scale(0); opacity: 0; }
    80%  { transform: scale(1.1); }
    100% { transform: scale(1); opacity: 1; }
  }

  .metrics {
    margin-top: 2rem;
    display: flex;
    gap: 2rem;
    justify-content: center;
    font-size: 0.8rem;
    color: #888;
  }
  .metric span {
    display: block;
    font-size: 1.4rem;
    color: #00d2ff;
    font-weight: bold;
  }
</style>
</head>
<body>

<div class="container">
  <h1>kube-watch</h1>
  <p style="color:#888; font-size:0.9rem;">kubernetes platform toolkit</p>

  <div class="status">
    <div class="dot"></div>
    Service healthy
  </div>

  <div class="pods">
    <div class="pod">pod/0</div>
    <div class="pod">pod/1</div>
    <div class="pod">pod/2</div>
    <div class="pod">pod/3</div>
  </div>

  <div class="metrics">
    <div class="metric">nodes<span id="nodes">0</span></div>
    <div class="metric">pods<span id="podCount">0</span></div>
    <div class="metric">cpu<span id="cpu">0%</span></div>
    <div class="metric">mem<span id="mem">0%</span></div>
  </div>
</div>

<script>
  // floating k8s wheels in background
  const svgWheel = '<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><circle cx="50" cy="50" r="44" fill="none" stroke="white" stroke-width="3"/><circle cx="50" cy="50" r="12" fill="white"/><line x1="50" y1="6" x2="50" y2="30" stroke="white" stroke-width="4" stroke-linecap="round"/><line x1="50" y1="70" x2="50" y2="94" stroke="white" stroke-width="4" stroke-linecap="round"/><line x1="6" y1="50" x2="30" y2="50" stroke="white" stroke-width="4" stroke-linecap="round"/><line x1="70" y1="50" x2="94" y2="50" stroke="white" stroke-width="4" stroke-linecap="round"/><line x1="19" y1="19" x2="36" y2="36" stroke="white" stroke-width="4" stroke-linecap="round"/><line x1="64" y1="64" x2="81" y2="81" stroke="white" stroke-width="4" stroke-linecap="round"/><line x1="81" y1="19" x2="64" y2="36" stroke="white" stroke-width="4" stroke-linecap="round"/><line x1="36" y1="64" x2="19" y2="81" stroke="white" stroke-width="4" stroke-linecap="round"/></svg>';

  for (let i = 0; i < 6; i++) {
    const el = document.createElement('div');
    el.className = 'wheel';
    el.innerHTML = svgWheel;
    const size = 60 + Math.random() * 120;
    el.style.width = size + 'px';
    el.style.height = size + 'px';
    el.style.left = Math.random() * 100 + 'vw';
    el.style.animationDuration = (15 + Math.random() * 20) + 's';
    el.style.animationDelay = -(Math.random() * 30) + 's';
    document.body.appendChild(el);
  }

  // animated counters
  function animateCount(id, target, suffix = '') {
    const el = document.getElementById(id);
    let current = 0;
    const step = Math.max(1, Math.floor(target / 30));
    const interval = setInterval(() => {
      current += step;
      if (current >= target) { current = target; clearInterval(interval); }
      el.textContent = current + suffix;
    }, 40);
  }
  setTimeout(() => animateCount('nodes', 4), 500);
  setTimeout(() => animateCount('podCount', 12), 700);
  setTimeout(() => animateCount('cpu', 37, '%'), 900);
  setTimeout(() => animateCount('mem', 64, '%'), 1100);
</script>

</body>
</html>`

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health" {
		fmt.Fprintf(w, "ok")
		return
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, page)
}

// The main fn, entry point of program exec, the main

func main(){
	// http.HandleFunc registers handler func to URL path : mapping root to healthCheckFn
	// "/" will be root path: localhost:8080/
	// when request are done to this path, we call healthCheckHandler
	http.HandleFunc("/", healthCheckHandler)

	// os.Getenv gets value from env named "PORT"
	// env variables are disjointed from main for security reasons

	port := os.Getenv("PORT")
	// := is for declaraing and assignment

	// we want to check if the port is empty:
	if port == "" {
		// set default to 8080
		port= "8080"
	}

	// log.printf( to print string to log/terminal)
	log.Printf("Starting Server on :%s", port)

	// http.ListenAndServe, starts http server, 2 arguments: 
	// 1. address to listen to
	// 2. handler fn ( nil to use fault set by HandleFunc)
	err := http.ListenAndServe(":"+port, nil)


	// if http.ListenAndServe returns error, block execution 
	if err != nil{
		// log.Fatalf error prints faral error message, and exits program w status 1
		log.Fatalf("Server Failed to Start: %v", err)
	}


}		