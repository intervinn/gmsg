import { useState } from "react"
import { useLocation } from "wouter"
import { api } from "../lib/api"

export default function AuthorizePage() {
    const [_, setLocation] = useLocation()
    const [action, setAction] = useState("register")
    const [error, setError] = useState("")

    const [formData, setFormData] = useState({} as any)

    const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const {name, value} = e.target
        setFormData((prev: any) => {
            return {
                ...prev,
                [name]: value
            }
        })
    }

    const register = () => {
        setError("")
        if (!formData.username || !formData.password) {
            setError("not all fields are filled")
            return
        }

        api.auth.register(formData.username, formData.password).then(res => {
            if (!res.ok) {
                setError(res.error)
                return
            }

            sessionStorage.setItem("user", JSON.stringify({
                id: res.data.id,
                username: res.data.username
            }))

            localStorage.setItem("token", "Bearer " + res.authorization)

            setLocation("/")
        })
    }

    const login = () => {
        setError("")
        if (!formData.username || !formData.password) {
            setError("not all fields are filled")
            return
        }

        api.auth.login(formData.username, formData.password).then(res => {
            if (!res.ok) {
                setError(res.error)
                return
            }

            sessionStorage.setItem("user", JSON.stringify({
                id: res.data.id,
                username: res.data.username
            }))

            localStorage.setItem("token", "Bearer " + res.authorization)

            setLocation("/")
        })
    }

    return <div className="h-screen w-full font-mono flex flex-col items-center justify-center">
        {action === "register" ? (
        <div className="flex flex-col">
            <h1 className="text-2xl">Register</h1>

            <form>
                <h2>Username</h2>
                <input onChange={onChange} value={formData.username || ""} name="username" type="text" className="border"/>
                <h2>Password</h2>
                <input onChange={onChange} value={formData.password || ""} name="password" type="text" className="border"/>
            </form>

            <button className="bg-slate-500 hover:bg-slate-600" onClick={register}>Register</button>
            <button className="text-blue-800" onClick={() => setAction("login")}>Login</button>
        </div>
    ) : (
        <div className="flex flex-col">
            <h1 className="text-2xl">Login</h1>

            <form>
                <h2>Username</h2>
                <input onChange={onChange} value={formData.username || ""} name="username" type="text" className="border"/>
                <h2>Password</h2>
                <input onChange={onChange} value={formData.password || ""} name="password" type="text" className="border"/>
            </form>

            <button className="bg-slate-500 hover:bg-slate-600" onClick={login}>Login</button>
            <button className="text-blue-800" onClick={() => setAction("register")}>Register</button>
        </div>
        )}

        <span className="text-red-500">{error}</span>
    </div>
}