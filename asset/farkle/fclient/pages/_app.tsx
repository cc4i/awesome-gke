import '../styles/globals.css'
import type { AppProps } from 'next/app'
import { createTheme, NextUIProvider } from '@nextui-org/react'
import { CssBaseline } from '@nextui-org/react';

export default function App({ Component, pageProps }: AppProps) {

  return (
    <NextUIProvider>
      {/* {CssBaseline.flush()} */}
      <Component {...pageProps} />
    </NextUIProvider>
    
  )
}
