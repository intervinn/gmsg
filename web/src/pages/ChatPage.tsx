import { Redirect } from "wouter";
import ChatWindow from "../components/ChatWindow";
import { useWebsocket, WebSocketProvider } from "../contexts/websocket";
import { useEffect } from "react";

export default function ChatPage() {
    const isAuthenticated = localStorage.getItem("token")

    if (!isAuthenticated) return <Redirect to="auth"/>

    const {lastMessage, sendMessage, isConnected} = useWebsocket()

    useEffect(() => {
        sendMessage({
            t: "authenticate",
            d: localStorage.getItem("token")
        })
    }, [])

    return (
        <div className="w-full h-screen font-mono">
            <WebSocketProvider>
                <ChatWindow/>
            </WebSocketProvider>
        </div>
    )
}