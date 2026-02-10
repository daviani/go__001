import { createSystem, defaultConfig, defineConfig} from "@chakra-ui/react";

const config = defineConfig({
    theme: {
        tokens: {
            colors: {
                nord: {
                    polar0: { value: '#2E3440' },
                    polar1: { value: '#3B4252' },
                    polar2: { value: '#434C5F' },
                    polar3: { value: '#4C566A' },
                    snow0: { value: '#D8DEE9' },
                    snow1: { value: '#E5E9F0'},
                    snow2: { value: '#ECEFF4'},
                    frost0:  { value: '#8FBCBB'},
                    frost1: { value: '#88C0D0'},
                    frost2: { value: '#81A1C1'},
                    frost3: { value: '#5E81AC'},
                    red: { value: '#BF616A'},
                    orange: { value: '#D08770'},
                    yellow: { value: '#EBCB8B'},
                    green: { value: '#A3BE8C'},
                    purple: { value: '#B48EAD' }
                }
            }
        },
        semanticTokens: {
            colors: {
                'bg.page':    { value: { base: '{colors.nord.snow2}', _dark: '{colors.nord.polar0}' } },
                'bg.card':    { value: { base: '{colors.nord.snow1}', _dark: '{colors.nord.polar1}' } },
                'text.main':  { value: { base: '{colors.nord.polar0}', _dark: '{colors.nord.snow1}' } },
                'text.muted': { value: { base: '{colors.nord.polar3}', _dark: '{colors.nord.snow0}' } },
                'accent':     { value: { base: '{colors.nord.frost2}', _dark: '{colors.nord.frost1}' } },
                'success':    { value: '{colors.nord.green}' },
                'error':      { value: '{colors.nord.red}' },
                'warning':    { value: '{colors.nord.orange}' },
            }
        }
    }
})

export const system = createSystem(defaultConfig, config)