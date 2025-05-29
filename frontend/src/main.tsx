import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';

const root = createRoot(document.getElementById('root')!);

// Function to render the application.
const renderApp = () => {
  root.render(
    <StrictMode>
      <App />
    </StrictMode>
  );
};

// Initial call to render the application.
renderApp();
