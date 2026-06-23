import { createContext, useContext, useEffect, useRef, useState, type RefObject } from "react";
import { kGatewayUrl } from "../lib/constants";

interface WebSocketContextProps {
    lastMessage: any,
    isConnected: boolean,
    sendMessage: (message: any) => void
}

const WebSocketContext = createContext<WebSocketContextProps | null>(null) 

export const WebSocketProvider = ({children} : {children : React.ReactNode}) => {
    const [lastMessage, setLastMessage] = useState(null as any)
    const [isConnected, setIsConnected] = useState(false)
    const sock : RefObject<null | WebSocket> = useRef(null)

    useEffect(() => {
        const socket = new WebSocket(kGatewayUrl)
        sock.current = socket

        socket.onopen = () => {
            setIsConnected(true)
        }

        socket.onclose = () => {
            setIsConnected(false)
        }

        socket.onmessage = (event) => {
            setLastMessage(JSON.parse(event.data))
        }

        return () => {
            socket.close()
        }
    }, [])

    const sendMessage = (message: any) => {
        if (sock.current && sock.current.readyState === WebSocket.OPEN) {
            sock.current.send(JSON.stringify(message))
        } else {
            console.error("websocket is not connected")
        }
    }

    return (
        <WebSocketContext.Provider value={{lastMessage, isConnected, sendMessage}}>
            {children}
        </WebSocketContext.Provider>
    )
}

export const useWebsocket = () => {
    const context = useContext(WebSocketContext)
    if (!context) {
        throw new Error("websocket context not found")
    }
    return context
}

