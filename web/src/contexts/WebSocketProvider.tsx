import { createContext, useRef, useState } from "react";

const WebSocketContext = createContext(null)

export const WebSocketProvider = ({children}) => {
    const [messages, setMessages] = useState([])
    const [isConnected, setIsConnected] = useState(false)
    const sock = useRef(null)

    
}