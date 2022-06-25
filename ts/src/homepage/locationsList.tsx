import {
    Tr,
    Th,
    Td,
    Tbody,
    Thead,
    Table,
    TableContainer,
    Tooltip,
    Center,
    VStack,
    Button,
    Popover,
    PopoverTrigger,
    PopoverContent,
    PopoverHeader,
    PopoverBody,
    PopoverFooter,
    PopoverArrow,
    PopoverCloseButton,
    ButtonGroup,
    useDisclosure,
    Input,
    useToast,
    HStack
} from '@chakra-ui/react';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { fetchLocations, addStockpile, deleteStockpile, refreshStockpile } from '../api/apis';

const LocationRow = ({ location, time, code, refreshFunc }: { location: string, time: Date, code: string, refreshFunc: () => void }) => {
    let expireTime = time.getTime() + 1000 * 3600 * 48;
    let hours = (expireTime - Date.now()) / (1000 * 3600);
    const navigate = useNavigate();

    return (
        <Tr key={location}>
            <Td>{location}</Td>
            <Td><Tooltip label={(new Date(expireTime)).toLocaleString()}>{Math.floor(hours).toString() + " hours"}</Tooltip></Td>
            <Td isNumeric>{code}</Td>
            <Td>
                <HStack>
                    <Button backgroundColor="green.500" size='sm'
                        onClick={() => { refreshStockpile(location, navigate); refreshFunc() }}>
                        Refresh
                    </Button>
                    <Button backgroundColor="red.500" size='sm'
                        onClick={() => { deleteStockpile(location, navigate); refreshFunc() }}>
                        Delete
                    </Button>
                </HStack>
            </Td>
        </Tr>
    );
};

const LocationsList = (props: any) => {
    const [rows, setRows] = useState<Array<JSX.Element>>([]);
    const [stockpileName, setStockpileName] = useState('');
    const [stockpileCode, setStockpileCode] = useState('');
    const [popoverButtonDisable, setPopoverButtonDisable] = useState(false);
    const navigate = useNavigate();
    const { onOpen, onClose, isOpen } = useDisclosure();
    const toast = useToast();
    const refreshRows = async () => {
        let locations = await fetchLocations();
        if (locations) {
            locations.sort((a, b) => {
                return a.time.valueOf() - b.time.valueOf();
            });
            setRows(locations.map(location => <LocationRow {...location} refreshFunc={refreshRows} />));
        }
    };

    useEffect(() => {
        refreshRows();
    }, []);

    const handlePopoverClose = () => {
        onClose();
    };
    const handlePopoverSubmit = () => {
        if (stockpileCode == '' || stockpileName == '') {
            toast({
                title: 'Failed to add stockpile',
                description: "Stockpile name and code mustn't be empty",
                status: 'error',
                duration: 9000,
                isClosable: true,
            });
            return;
        }
        setPopoverButtonDisable(true);
        addStockpile(stockpileName, stockpileCode, navigate).then(success => {
            if (!success) {
                toast({
                    title: 'Failed to add stockpile',
                    description: "The stockpile already exists",
                    status: 'error',
                    duration: 9000,
                    isClosable: true,
                });
            } else {
                toast({
                    title: 'Stockpile added',
                    description: "The stockpile has been added",
                    status: 'success',
                    duration: 9000,
                    isClosable: true,
                });
                refreshRows();
            }
            setPopoverButtonDisable(false);
            onClose();
        });
    };

    return (
        <Center>
            <VStack width={[
                '100%', // 0-30em
                '70%', // 30em-48em
                '60%', // 48em-62em
            ]}>
                <Popover
                    placement='bottom'
                    closeOnBlur={false}
                    isOpen={isOpen}
                >
                    <PopoverTrigger>
                        <Button alignSelf="flex-end" onClick={() => onOpen()}>Add Stockpile</Button>
                    </PopoverTrigger>
                    <PopoverContent >
                        <PopoverHeader fontWeight='bold'>
                            Add a new stockpile
                        </PopoverHeader>
                        <PopoverArrow />
                        <PopoverCloseButton onClick={() => onClose()}/>
                        <PopoverBody>
                            <VStack>
                                <Input placeholder='Stockpile Location' onChange={event => setStockpileName(event.target.value)} />
                                <Input placeholder='Code' onChange={event => setStockpileCode(event.target.value)} />
                            </VStack>
                        </PopoverBody>
                        <PopoverFooter
                            display='flex'
                            justifyContent='flex-end'
                        >
                            <ButtonGroup size='sm'>
                                <Button colorScheme='green'
                                    disabled={popoverButtonDisable}
                                    onClick={handlePopoverSubmit}>
                                    Comfirm
                                </Button>
                                <Button colorScheme='blue'
                                    onClick={handlePopoverClose}>Close</Button>
                            </ButtonGroup>
                        </PopoverFooter>
                    </PopoverContent>
                </Popover>
                <TableContainer width="100%">
                    <Table variant='simple'>
                        <Thead>
                            <Tr>
                                <Th>Stockpile Name</Th>
                                <Th>Expire In</Th>
                                <Th isNumeric>Code</Th>
                                <Th></Th>
                            </Tr>
                        </Thead>
                        <Tbody>
                            {rows}
                        </Tbody>
                    </Table>
                </TableContainer>
            </VStack>
        </Center>
    );
};

export { LocationsList };