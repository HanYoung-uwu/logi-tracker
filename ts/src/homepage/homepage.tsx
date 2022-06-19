import { Text, Tabs, TabList, TabPanels, Tab, TabPanel } from '@chakra-ui/react';
import { observer } from 'mobx-react-lite';
import { Navigate } from 'react-router-dom';
import { ItemsList } from './itemsList';
import { HistoryList } from './historyList';
import { LocationsList } from './locationsList';
import { useContext } from 'react';
import { AccountInfoContext } from '../login/accountInfoStore';

const HomePage = observer((props: any) => {
    const accountInfo = useContext(AccountInfoContext);

    if (accountInfo.permission != -1 && accountInfo.permission != 0 && accountInfo.permission != 1 && accountInfo.permission != 2) {
        return <Navigate to="/login" replace={true} />;
    }

    return (
        <Tabs isLazy>
            <TabList>
                <Tab marginLeft={["0", "10%", "20%"]}>Items</Tab>
                <Tab>Stockpiles</Tab>
                <Tab>History</Tab>
            </TabList>

            <TabPanels>
                <TabPanel>
                    <ItemsList />
                </TabPanel>
                <TabPanel>
                    <LocationsList />
                </TabPanel>
                <TabPanel>
                    <HistoryList />
                </TabPanel>
            </TabPanels>
        </Tabs>
    );
});

export { HomePage };