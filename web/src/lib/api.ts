import { kApiUrl } from "./constants"

export type Response<T> = {
    ok: boolean,
    message?: string,
    error?: string,
    data?: T
}

export const api = {
    auth: {
        register: async (username: string, password: string) => {
            const res = await fetch(kApiUrl + "/auth/register", {
                method: "POST",
                body: JSON.stringify({username, password})
            })

            return await res.json()
        },

        login: async (username: string, password: string) => {
            const res = await fetch(kApiUrl + "/auth/login", {
                method: "POST",
                body: JSON.stringify({username, password})
            })

            return await res.json()
        }
    }
}