import { HStack, VStack, Image, Text, Input, Button } from "@chakra-ui/react";
import { useState, useEffect } from "react";
import { IconList } from "../icons/icons";
import { MAX_ITEM_TABLE_DISPLAY_SIZE } from "../config/config";

const ItemSelectPage = (props: any) => {
    const [filter, setFilter] = useState('');

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


    return (
        <VStack>
            <Input onChange={event => setFilter(event.target.value)}
                placeholder="Search..."
                width={[
                    '100%', // 0-30em
                    '80%', // 30em-48em
                    '50%', // 48em-62em
                ]} />
            <ItemsTable filter={filter} onClick={props.onClick} rowNum={rowNumber} />
        </VStack>)
};

const PoorMansPagination = (props: { size: number, onClick: (page: number) => void, currentPage: number }) => {
    let buttonSize = props.size;
    let buttons = [];
    for (let i = 0; i < buttonSize; i++) {
        buttons.push(<Button onClick={() => props.onClick(i)}
            disabled={i == props.currentPage ? true : false}>
            {(i + 1).toString()}
        </Button>);
    };
    return (
        <HStack>
            {buttons}
        </HStack>
    )
};

const ItemsTable = (props: { filter?: string, rowNum?: number, onClick?: (name: string) => any }) => {
    const [currentPage, setCurrentPage] = useState(0);

    let sources = IconList;

    if (props.filter && props.filter != "") {
        let filter: string = props.filter.toLowerCase();
        sources = sources.filter(({ name }) => name.toLowerCase().includes(filter));
    }
    let rowNum = 8;
    if (props.rowNum) {
        rowNum = props.rowNum;
    }

    let imgs = [];
    let row: JSX.Element[] = [];
    for (let i = MAX_ITEM_TABLE_DISPLAY_SIZE * currentPage; i < sources.length; i++) {
        if (i >= MAX_ITEM_TABLE_DISPLAY_SIZE * (currentPage + 1)) {
            break;
        };

        if (row.length == rowNum) {
            imgs.push(<HStack>{row}</HStack>);
            row = [];
        }
        row.push(<VStack key={sources[i].name}
            onClick={() => { if (props.onClick) props.onClick(sources[i].name); }}>
            <Image
                boxSize='100px'
                objectFit='contain'
                src={sources[i].icon}
                alt={sources[i].name}
            />
            <Text maxW="100px">{sources[i].name}</Text>
        </VStack>);
    }

    if (row.length != 0) {
        imgs.push(<HStack>{row}</HStack>);
    }
    
    let buttonSize = Math.ceil(sources.length / MAX_ITEM_TABLE_DISPLAY_SIZE);
    imgs.push(<PoorMansPagination onClick={(page: number) => setCurrentPage(page)}
        size={buttonSize}
        currentPage={currentPage} />);

    return (<VStack>
        {imgs}
    </VStack>);
}

export default ItemSelectPage;