import { useState } from "react"

interface Message {
    content: string,
    author: {
        username: string
    }
}

export default function ChatWindow() {
    const [messages, setMessages] = useState<Message[]>([])
    
    const [input, setInput] = useState("")

    return (
        <div className="w-full h-full grid grid-rows-[1fr_auto] gap-2">
            <div className="flex-col">
                {messages.map(msg => 
                    <div className="px-3 py-2">
                        <h3>{msg.author.username}</h3>
                        <span className="">{msg.content}</span>
                    </div>
                )}
            </div>
            <div className="flex-row">
                <input type="text" placeholder="Type something here..." />
                <button>Send</button>
            </div>
        </div>
    )
}