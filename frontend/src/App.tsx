import { BrowserRouter, Route, Routes } from "react-router-dom"
import LoginRegister from "./pages/LoginRegister/LoginRegister"
import Main from "./pages/main/Main"


function App() {
 

  return (
    <>
      <BrowserRouter>
        <Routes>
          
          <Route path="/" element={<LoginRegister />} />
          <Route path="/main" element={<Main />} />
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
