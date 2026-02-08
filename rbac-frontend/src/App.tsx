import {
  BrowserRouter as Router,
  Routes,
  Route,
  useLocation,
} from "react-router-dom";
import { Dashboard } from "./pages/Dashboard";
import Login from "./pages/Login";
import { CreateUser } from "./pages/CreateUser";
import { CreateProject } from "./pages/CreateProject";
import Tasks from "./pages/Tasks";
import { Navbar } from "./components/Navbar";
import ProjectDetails from "./pages/ProjectDetails";

function App() {
  return (
    <Router>
      <Inner />
    </Router>
  );
}

function Inner() {
  const location = useLocation();
  const hideNav = location.pathname === "/";
  return (
    <div>
      {!hideNav && <Navbar />}
      <Routes>
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/create-user" element={<CreateUser />} />
        <Route path="/create-project" element={<CreateProject />} />
        <Route path="/tasks" element={<Tasks />} />
        <Route path="/projects/:id" element={<ProjectDetails />} />
        <Route path="/" element={<Login />} />
      </Routes>
    </div>
  );
}

export default App;