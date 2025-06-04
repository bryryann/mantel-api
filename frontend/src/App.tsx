import { useEffect } from 'react';
import { proxyHealthcheck } from './services/proxy-service';
import AppRoutes from '@/routes/AppRoutes';
import './main.css';
import axios from 'axios';

const App = () => {
  useEffect(() => {
    proxyHealthcheck();
  });

  return <AppRoutes />;
};

export default App;
