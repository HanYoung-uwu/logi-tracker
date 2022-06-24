import { HStack, VStack, Image, Text, Tabs, TabList, TabPanels, Tab, TabPanel } from "@chakra-ui/react";
import { useState, useEffect } from "react";
import { IconList } from "../icons/icons";

const ItemSelectPage = (props: any) => {
    const [rowNumber, setRowNumber] = useState(8);
    useEffect(() => {
        function handleResize() {
            if (window.innerWidth >= 1080) {
                setRowNumber(8);
            } else if (window.innerWidth >= 720) {
                setRowNumber(5);
            } else {
                setRowNumber(3);
            }
        }
        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, []);


    return (<ItemsTable onClick={props.onClick} rowNum={rowNumber} />)
};

const ItemsTable = (props: { rowNum: number, onClick?: (name: string) => any }) => {
    let sources = IconList;
    let rowNum = props.rowNum;
    const IconWithDesc = ({ name, icon }: { name: string, icon: string }) => {
        return (<VStack key={name}
            onClick={() => { if (props.onClick) props.onClick(name); }}>
            <Image
                boxSize='100px'
                objectFit='contain'
                src={icon}
                alt={name}
            />
            <Text maxW="100px">{name}</Text>
        </VStack>);
    };

    const IconsTable = ({ icons }: { icons: Array<{ name: string, icon: string }> }) => {
        let result = [];
        let row: Array<JSX.Element> = [];
        for (let i = 0; i < icons.length; i++) {
            if (row.length == rowNum) {
                result.push(<HStack key={i}>{row}</HStack>);
                row = [];
            }
            row.push(<IconWithDesc {...icons[i]} />);
        }
        if (row.length != 0) {
            result.push(<HStack key={icons.length}>{row}</HStack>);
        }
        return (<VStack>{result}</VStack>)
    };

    return (<Tabs isLazy>
        <TabList>
            <Tab>Small Arms</Tab>
            <Tab>Heavy Arms</Tab>
            <Tab>Heavy Ammunition</Tab>
            <Tab>Utility</Tab>
            <Tab>Supply</Tab>
            <Tab>Medical</Tab>
            <Tab>Uniforms</Tab>
            <Tab>Resource</Tab>
            <Tab>Vehicle</Tab>
            <Tab>Shippable</Tab>
        </TabList>

        <TabPanels>
            <TabPanel>
                <IconsTable icons={sources["small arms"]} />
            </TabPanel>
            <TabPanel>
                <IconsTable icons={sources["heavy arms"]} />
            </TabPanel>
            <TabPanel>
                <IconsTable icons={sources["heavy ammunition"]} />
            </TabPanel>
            <TabPanel>
                <IconsTable icons={sources.tool} />
            </TabPanel>
            <TabPanel>
                <IconsTable icons={sources.supply} />
            </TabPanel>
            <TabPanel>
                <IconsTable icons={sources.medical} />
            </TabPanel>
            <TabPanel>
                <IconsTable icons={sources.uniform} />
            </TabPanel>
            <TabPanel>
                <IconsTable icons={sources.resource} />
            </TabPanel>
            <TabPanel>
                <IconsTable icons={sources.vehicle} />
            </TabPanel>
            <TabPanel>
                <IconsTable icons={sources.shippable} />
            </TabPanel>
        </TabPanels>
    </Tabs>);
}

export default ItemSelectPage;