import { useEffect } from "react";

const App = () => {
  console.log(import.meta.env.VITE_API_URL);

  useEffect(() => {
    fetch('/api/healthcheck')
      .then(res => res.json())
      .then(data => console.log('API response: ', data))
      .catch(err => console.error('API error: ', err));
  }, []);

  return <h1>mantel.</h1>;
};

export default App;
