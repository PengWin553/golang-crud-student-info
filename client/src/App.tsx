import React from 'react';
// import Navbar from "./components/Navbar";
// import TodoForm from "./components/TodoForm";
import Students from "./components/Students";

export const BASE_URL = import.meta.env.MODE === "development" ? "http://localhost:4000/api" : "/api";

function App() {
  return (
    <div>
      {/* <Navbar /> */}
      <div>
        {/* <TodoForm /> */}
        <Students /> 
      </div>
    </div>
  );
}

export default App;