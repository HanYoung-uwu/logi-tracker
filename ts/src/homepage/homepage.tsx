import {
    Input, Flex, Text, Table,
    Center,
    Thead,
    Tbody,
    Tfoot,
    Select,
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
} from '@chakra-ui/react'
import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom';
import { addItem, fetchAllItems, fetchLocations, Location } from '../api/apis'
import ItemSelectPage from '../itemTable/itemTable';


const ItemsTable = (props: {fetchRef: ((arg0: Function) => any)} | null) => {
    const [rows, setRows] = useState(Array<JSX.Element>());
    let navigate = useNavigate();
    const fetchAndConstructTable = async () => {
        let items = await fetchAllItems(navigate);
        if (items) {
            setRows(items.map((record) => {
                return (<Tr key={record.item + "###" + record.location}>
                    <Td>{record.item}</Td>
                    <Td isNumeric>{record.size}</Td>
                    <Td>{record.location}</Td>
                </Tr>)
            }));
        }
    }

    props?.fetchRef(fetchAndConstructTable);

    useEffect(() => {
        fetchAndConstructTable();
    }, []);

    return (
        <TableContainer width="100%">
            <Table variant='simple'>
                <Thead>
                    <Tr>
                        <Th>Item</Th>
                        <Th isNumeric>Size</Th>
                        <Th>Location</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    {rows}
                </Tbody>
            </Table>
        </TableContainer>);
}

const HomePage = (props: any) => {
    const { isOpen, onOpen, onClose } = useDisclosure();
    const [locations, setLocations] = useState<Array<Location>>();
    const [selectedLocation, setSelectedLocation] = useState<string>('');
    const [quantity, setQuantity] = useState<number>(0);
    
    let refreshTable: any;

    const handleAddItem = (name: string) => {
        let location = selectedLocation;
        if (location == '' && locations) {
            location = locations[0].location;
        }
        addItem(name, quantity, location);
        onClose();
        refreshTable();
    };
    useEffect(() => {
        const initLocation = async () => {
            let res = await fetchLocations();
            if (res != null) {
                setLocations(res);
            }
        };
        initLocation();
    }, []);

    return (
        <Center>
            <VStack width={[
                '100%', // 0-30em
                '70%', // 30em-48em
                '60%', // 48em-62em
            ]}>
                <Button alignSelf="flex-end" onClick={() => onOpen()}>Add Item</Button>
                <ItemsTable fetchRef={refresh => refreshTable = refresh}/>
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
                                        {locations?.map(({ location }) => <option value={location}>{location}</option>)}
                                    </Select>
                                    <Input onChange={event => setQuantity(Number.parseInt(event.target.value))} type="number" placeholder='quantity'></Input>
                                </HStack>
                                <Text>Click on item to add</Text>
                            <ItemSelectPage onClick={handleAddItem}/>
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
}

export { HomePage }