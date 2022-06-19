import * as React from "react";
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
} from "@chakra-ui/react";

import Theme from "./theme";
import LoginPage from "./login/loginpage";
import { HomePage } from "./homepage/homepage";
import HeaderPanel from "./navigationHeader/navigationHeader";
import RegisterPage from "./registerPage/registerPage";
import ItemsPage from "./itemTable/itemTable";
import { AccountInfoContext, AccountInfo } from "./login/accountInfoStore";

export const App = () => (
  <BrowserRouter>
    <ChakraProvider theme={Theme}>
      <AccountInfoContext.Provider value={new AccountInfo()}>
        <HeaderPanel />
        <Routes>
          <Route path="/login" element={<LoginPage width={[
            '100%', // 0-30em
            '50%', // 30em-48em
            '25%', // 48em-62em
          ]} />} />
          <Route path="/home" element={<HomePage />} />
          <Route path="/" element={<Text>Root page</Text>} />
          <Route path="/invite" element={<RegisterPage />} />
          <Route path="/items" element={<ItemsPage />} />
        </Routes>
      </AccountInfoContext.Provider>
    </ChakraProvider>
  </BrowserRouter>
)
