import {Button, Flex, Input, NativeSelect, Text} from "@chakra-ui/react";
import {useState} from "react";
import {isValidDomain, isValidExtension} from "../utils/validation.ts";
import {toaster} from "./ui/toaster.tsx";

function ScanForm({ onScan }: { onScan: (domain: string, scanType: string) => void }) {
    const [domain, setDomain] = useState("")
    const [scanType, setScanType] = useState("all")

    const handleClick = () => {
        const trimmed = domain.trim()
        if (!isValidDomain(trimmed)) {
            toaster.create({ title: "Format de domaine invalide", type: "warning" })
            return
        }
        if (!isValidExtension(trimmed)) {
            toaster.create({ title: "Extension de domaine inconnue", type: "warning" })
            return
        }
        onScan(trimmed, scanType)
    }

    return (
        <Flex p="8" maxW="600px" mx="auto" direction="column" gap="4" >
            <Text p="8" textAlign="center" >Domaine à analyser</Text>
            <Input  bg="bg.card"
                    borderColor="nord.polar3"
                    color="text.main"
                    value={domain}
                    onChange={(e) => setDomain(e.target.value)}
                    placeholder="daviani.dev"
            />
            <NativeSelect.Root>
                <NativeSelect.Field
                    value={scanType}
                    onChange={(e) => setScanType(e.target.value)}
                    borderColor="nord.polar3" color="text.main"
                >
                    <option value="all">Tous les scanners</option>
                    <option value="dns">DNS</option>
                    <option value="ssl">SSL/TLS</option>
                    <option value="header">Headers HTTP</option>
                    <option value="subdomain">Sous-domaines</option>
                    <option value="sensitive">Fichiers sensibles</option>
                </NativeSelect.Field>
            </NativeSelect.Root>
            <Button disabled={domain.trim() === ""} bg="accent" color="nord.polar0" onClick={handleClick}>Scanner</Button>
        </Flex>
    )
}

export default ScanForm;
