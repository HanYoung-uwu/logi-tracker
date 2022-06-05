import {
    Input, Flex, Text, Table,
    Thead,
    Tbody,
    Tfoot,
    Tr,
    Th,
    Td,
    TableCaption,
    TableContainer,
} from '@chakra-ui/react'
import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom';
import { fetchAllItems } from '../api/apis'


const ItemsTable = (props: any) => {
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

    useEffect(() => {
        fetchAndConstructTable();
    }, []);

    return (<TableContainer>
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
    </TableContainer>)
}

export { ItemsTable }