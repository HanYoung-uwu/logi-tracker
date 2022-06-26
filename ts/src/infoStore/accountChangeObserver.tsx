import { useContext, useEffect } from "react";
import { observer } from "mobx-react-lite";
import { AccountInfoContext, AccountInfo } from "./accountInfoStore";
import { LocationInfoContext, StockpileInfo } from "./stockPileInfoStore";

const AccountChangeObserver = observer(() => {
    const accountInfo = useContext(AccountInfoContext);
    const stockpileInfo = useContext(LocationInfoContext);

    useEffect(() => {
        if (accountInfo.name !== "") {
            stockpileInfo.refetchInfo();
        }
    }, [accountInfo.name]);
    return null;
});

export default AccountChangeObserver;