import Header from "./components/Header.tsx";
import ScanForm from "./components/ScanForm.tsx";
import {Toaster, toaster} from "./components/ui/toaster.tsx";
import {useState} from "react";
import ScanResults from "./components/ScanResults.tsx";
import {scanDomain, type ScanResult} from "./services/scanner.ts";
import {Box} from "@chakra-ui/react";

function App() {
    const [results, setResults] = useState<ScanResult[]>([])
    const [loading, setLoading] = useState(false)

    const handleScan = async (
        domain:string,
        scanType:string
    ) => {
        setLoading(true)
        try {
            const data = await scanDomain(domain, scanType);
            setResults(data);
        } catch (e) {
            toaster.create({ title: "Erreur lors du scan", type: "error" })
        } finally {
            setLoading(false)
        }
    }
    return (
        <Box bg="bg.page" minH="100vh">
            <Header />
            <ScanForm onScan={handleScan} />
            <ScanResults results={results} loading={loading} />
            <Toaster />
        </Box>
    )
}

export default App
