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
import { AccountInfoContext, AccountInfo } from "./infoStore/accountInfoStore";
import { LocationInfoContext, StockpileInfo } from "./infoStore/stockPileInfoStore";
import AccountChangeObserver from "./infoStore/accountChangeObserver";
import RootPage from "./rootPage/rootPage";
import ClanManagePage from "./manageClan/clanManagePage";

export const App = () => (
  <BrowserRouter>
    <ChakraProvider theme={Theme}>
      <AccountInfoContext.Provider value={new AccountInfo()}>
        <LocationInfoContext.Provider value={new StockpileInfo()}>
          <AccountChangeObserver />
          <HeaderPanel />
          <Routes>
            <Route path="/login" element={<LoginPage width={[
              '100%', // 0-30em
              '50%', // 30em-48em
              '25%', // 48em-62em
            ]} />} />
            <Route path="/home" element={<HomePage />} />
            <Route path="/" element={<RootPage />} />
            <Route path="/invite" element={<RegisterPage />} />
            <Route path="/items" element={<ItemsPage />} />
            <Route path="/manage" element={<ClanManagePage />} />
          </Routes>
        </LocationInfoContext.Provider>
      </AccountInfoContext.Provider>
    </ChakraProvider>
  </BrowserRouter>
)
