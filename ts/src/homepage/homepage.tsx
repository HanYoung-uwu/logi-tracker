import { Text, Tabs, TabList, TabPanels, Tab, TabPanel } from '@chakra-ui/react'
import { ItemsList } from './itemsList';
import { HistoryList } from './historyList';
import { LocationsList } from './locationsList';

const HomePage = (props: any) => {
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
};

export { HomePage };