import {
    ListItem,
    Center,
    UnorderedList,
} from '@chakra-ui/react'
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { fetchHistory, HistoryRecord } from '../api/apis';

const HistoryList = (props: any) => {
    const [listSource, setListSource] = useState<Array<JSX.Element>>([]);
    const navigate = useNavigate();
    useEffect(() => {
        const initList = async () => {
            let ret = await fetchHistory(navigate);
            let list = ret.map(({ User, Location, ItemType, Size, Time, Action }) => {
                let desc = "";
                switch (Action) {
                    case "add":
                        desc = `${User} ${Action}ed ${Size} ${ItemType} into ${Location} - ${Time.toLocaleDateString()}`; break;
                    case "retrieve":
                        desc = `${User} ${Action}d ${Size} of ${ItemType} from ${Location} - ${Time.toLocaleDateString()}`; break;
                    case "delete":
                        desc = `${User} ${Action}d ${ItemType} at ${Location} - ${Time.toLocaleDateString()}`; break;
                    case "set":
                        desc = `${User} ${Action} the quantity of ${ItemType} at ${Location} to ${Size} - ${Time.toLocaleDateString()}`; break;
                    case "create stockpile":
                        desc = `${User} ${Action} ${Location} - ${Time.toLocaleDateString()}`; break;
                    case "delete stockpile":
                        desc = `${User} ${Action} ${Location} - ${Time.toLocaleDateString()}`; break;
                }
                return (
                    <ListItem key={User + Location + ItemType + Size + Time.toString() + Action}>
                        {desc}
                    </ListItem>
                )
            });

            setListSource(list);
        };
        initList();
    }, []);

    return (
        <Center>
            <UnorderedList width={[
                '100%', // 0-30em
                '70%', // 30em-48em
                '60%', // 48em-62em
            ]}>
                {listSource}
            </UnorderedList>
        </Center>);
}

export { HistoryList };