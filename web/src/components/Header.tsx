import {Button, Flex, Text} from "@chakra-ui/react";
import {useColorMode} from "./ui/color-mode.tsx";

function Header() {
    const {colorMode, toggleColorMode } = useColorMode()
    return (
        <>
            <Flex justify="space-between" align="center" p="4" bg="bg.header">
                <Text color="nord.snow1" textStyle="2xl" fontWeight="bold" >
                    Go Sentry - Scan
                </Text>
                <Button onClick={toggleColorMode} bg="bg.header" >
                    {colorMode === 'light' ? '🌙' : '☀️'}
                </Button>
            </Flex>
        </>
    )
}

export default Header;