const API_URL = "http://localhost:8082"

export interface ScanResult {
    scanner: string
    domain: string
    result: string
}

export async function scanDomain(domain: string, scanType: string): Promise<ScanResult[]> {
    const res = await fetch(`${API_URL}/scan/${scanType}?domain=${domain}`)
    if (!res.ok) {
        throw new Error('Erreur serveur')
    }
    const data = await res.json()
    return Array.isArray(data) ? data : [data]
}