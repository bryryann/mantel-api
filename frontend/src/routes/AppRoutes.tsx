import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import AuthorizationScreen from "@/pages/AuthorizationScreen";
import NotFound from "@/pages/NotFound/NotFound";

const AppRoutes = () => {
  return (
    <Router >
      <Routes>
        <Route path='/' element={<AuthorizationScreen />} />

        {/* Catch-all route for 404 */}
        <Route path='*' element={<NotFound />} />
      </Routes>
    </Router>
  );
};

export default AppRoutes;
