import React from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'

const root = createRoot(document.getElementById('root')!)

function renderApplication() {
  root.render(
    <React.StrictMode>
      <App />
    </React.StrictMode>
  );
}

renderApplication();
