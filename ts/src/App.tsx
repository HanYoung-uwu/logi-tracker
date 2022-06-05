import * as React from "react"
import { BrowserRouter, Routes, Route, Link as RLink } from "react-router-dom";

import {
  ChakraProvider,
  Box,
  Text,
  Link,
  VStack,
  Code,
  Grid,
  Flex
} from "@chakra-ui/react"

import Theme from "./theme"
import LoginPage from "./login/loginpage"
import { ItemsTable } from "./allitems.tsx/allitemstable";
export const App = () => (
  <BrowserRouter>
    <ChakraProvider theme={Theme}>
      <nav>
        <Link as={RLink} to="/login">login</Link>
      </nav>
      <Routes>
        <Route path="/login" element={<LoginPage width={[
          '100%', // 0-30em
          '50%', // 30em-48em
          '25%', // 48em-62em
        ]} />} />
        <Route path="/home" element={<ItemsTable />} />
        <Route path="/" element={<Text>Root page</Text>} />
      </Routes>
    </ChakraProvider>
  </BrowserRouter>
)
