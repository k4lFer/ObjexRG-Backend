import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom"
import { Toaster } from "sonner"

import { AppLayout } from "./components/AppLayout"
import { ProtectedRoute } from "./components/ProtectedRoute"
import { AuthPage } from "./pages/auth/AuthPage"
import { AuthProvider } from "./pages/auth/application/AuthProvider"
import { DrivePage } from "./pages/drive/DrivePage"
import { LandingPage } from "./pages/landing/LandingPage"
import { UserPage } from "./pages/user/UserPage"

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/login" element={<AuthPage />} />
          <Route element={<ProtectedRoute />}>
            <Route element={<AppLayout />}>
              <Route path="/drive" element={<DrivePage />} />
              <Route path="/profile" element={<UserPage />} />
            </Route>
          </Route>
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
        <Toaster richColors position="top-right" />
      </AuthProvider>
    </BrowserRouter>
  )
}

export default App
