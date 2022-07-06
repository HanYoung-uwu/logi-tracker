import {
    Input,
    Text,
    Table,
    Center,
    Thead,
    Tbody,
    Select,
    Editable,
    EditablePreview,
    EditableInput,
    Tr,
    Th,
    Td,
    TableContainer,
    VStack,
    Button,
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
    useDisclosure,
    HStack,
    Popover,
    PopoverTrigger,
    PopoverContent,
    PopoverHeader,
    PopoverBody,
    PopoverFooter,
    PopoverArrow,
    PopoverCloseButton,
    Menu,
    MenuButton,
    MenuList,
    MenuItem,
    Icon
} from '@chakra-ui/react';
import { observer } from 'mobx-react-lite';
import { useState, useContext } from 'react';
import { FaChevronDown } from 'react-icons/fa';
import { addItem, deleteItem, setItem } from '../api/apis'
import ItemSelectPage from '../itemTable/itemTable';
import { LocationInfoContext, StockpileInfo } from "../infoStore/stockPileInfoStore";
import { ItemsStoreContext } from '../infoStore/itemsStore';

const ItemTableRow = (props: { record: { location: string, size: number, item: string }, refreshCallback: () => any | null }) => {
    let record = props.record;
    const { onOpen, onClose, isOpen } = useDisclosure();
    const handleDeleteItem = () => {
        deleteItem(record.item, record.location);
        props.refreshCallback();
    };
    const [numberSize, setNumberSize] = useState(record.size.toString());

    const handleSetItem = () => {
        let tmp = Number.parseInt(numberSize);
        if (tmp != NaN) {
            setItem(record.item, tmp, record.location);
        } else {
            setNumberSize(record.size.toString());
        }
        onClose();
    };
    const handleChangeInput = (newVal: string) => {
        let tmp = Number.parseInt(newVal);
        if (!isNaN(tmp)) {
            setNumberSize(newVal);
        }
    };

    const handleSubmit = (val: string) => {
        if (Number.parseInt(val) != record.size) {
            onOpen();
        }
    };

    const handleClosePopover = () => {
        onClose();
        setNumberSize(record.size.toString());
    };

    return (<Tr key={record.item + "###" + record.location}>
        <Td>{record.item}</Td>
        <Td isNumeric>
            <Popover isOpen={isOpen}>
                <PopoverTrigger>
                    <Editable onSubmit={handleSubmit}
                        value={numberSize.toString()}
                        onChange={handleChangeInput}>
                        <EditablePreview />
                        <EditableInput type="number" />
                    </Editable>
                </PopoverTrigger>
                <PopoverContent>
                    <PopoverArrow />
                    <PopoverCloseButton onClick={handleClosePopover} />
                    <PopoverHeader alignSelf="flex-start">Comfirmation</PopoverHeader>
                    <PopoverBody>Are you sure you want to set the quantity?</PopoverBody>
                    <PopoverFooter>
                        <Button size='sm' colorScheme='green' onClick={handleSetItem}>Yes</Button>
                    </PopoverFooter>
                </PopoverContent>
            </Popover>
        </Td>
        <Td>{record.location}</Td>
        <Td><Button backgroundColor="red.500" onClick={handleDeleteItem}>Delete</Button></Td>
    </Tr>)
};

const ItemsTable = observer(({ filterLocation, fetchRef }: { filterLocation: string, fetchRef: ((arg0: Function) => any) }) => {
    const itemsStore = useContext(ItemsStoreContext);

    let filteredRows = itemsStore.getItems();
    if (filterLocation !== "") {
        filteredRows = itemsStore.getItems().filter(({ location }) => location === filterLocation);
    }

    let displayedRows = filteredRows.map(record => <ItemTableRow record={record}
        key={JSON.stringify(record)}
        refreshCallback={() => itemsStore.refetchInfo()} />);

    return (
        <TableContainer width="100%">
            <Table variant='simple'>
                <Thead>
                    <Tr>
                        <Th>Item</Th>
                        <Th isNumeric>Size</Th>
                        <Th>Location</Th>
                        <Th></Th>
                    </Tr>
                </Thead>
                <Tbody>
                    {displayedRows}
                </Tbody>
            </Table>
        </TableContainer>);
});

const ItemsList = observer((props: any) => {
    const { isOpen, onOpen, onClose } = useDisclosure();
    const [selectedLocation, setSelectedLocation] = useState<string>('');
    const [filterLocation, setFilterLocation] = useState<string>('');
    const [quantity, setQuantity] = useState<number>(0);

    const stockpileInfo = useContext(LocationInfoContext);

    let locations = stockpileInfo.locations;

    let refreshTable: any;

    const handleAddItem = (name: string) => {
        let location = selectedLocation;
        if (location == '' && locations.length !== 0) {
            location = locations[0].location;
        }
        if (location.length !== 0 && quantity !== 0) {
            addItem(name, quantity, location);
            onClose();
            refreshTable();
        }
    };

    const menuItems = () => {
        let result = [
            <MenuItem
                onClick={() => setFilterLocation("")}
                key="all">
                All
            </MenuItem>
        ];
        for (let i = 0; i < locations.length; i++) {
            let { location } = locations[i];
            result.push(
                <MenuItem
                    onClick={() => setFilterLocation(location)}
                    key={location}>
                    {location}
                </MenuItem>);
        }
        return result;
    };

    return (
        <Center>
            <VStack width={[
                '100%', // 0-30em
                '70%', // 30em-48em
                '60%', // 48em-62em
            ]}>

                <HStack alignSelf="flex-end">
                    <Menu>
                        <MenuButton as={Button} rightIcon={<Icon as={FaChevronDown} />}>
                            {filterLocation === "" ? "Location" : filterLocation}
                        </MenuButton>
                        <MenuList>
                            {menuItems()}
                        </MenuList>
                    </Menu>
                    <Button onClick={() => onOpen()}>Add Item</Button>
                </HStack>
                <ItemsTable fetchRef={refresh => refreshTable = refresh} filterLocation={filterLocation} />
            </VStack>
            <Modal isOpen={isOpen} onClose={onClose} size='full'>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>Add Item</ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <VStack>
                            <HStack>
                                <Text>Stockpile</Text>
                                <Select onChange={event => setSelectedLocation(event.target.value)}>
                                    {locations?.map(({ location }) => <option key={location} value={location}>{location}</option>)}
                                </Select>
                                <Input onChange={event => setQuantity(Number.parseInt(event.target.value))} type="number" placeholder='quantity'></Input>
                            </HStack>
                            <Text>Click on item to add</Text>
                            <ItemSelectPage onClick={handleAddItem} />
                        </VStack>
                    </ModalBody>
                    <ModalFooter>
                        <Button colorScheme='blue' mr={3} onClick={onClose}>
                            Close
                        </Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </Center>
    );
});

export { ItemsList }