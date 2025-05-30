import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import SplashScreen from '@/pages/SplashScreen';
import NotFound from "@/pages/NotFound/NotFound";

const AppRoutes = () => {
  return (
    <Router >
      <Routes>
        <Route path='/' element={<SplashScreen />} />

        {/* Catch-all route for 404 */}
        <Route path='*' element={<NotFound />} />
      </Routes>
    </Router>
  );
};

export default AppRoutes;
