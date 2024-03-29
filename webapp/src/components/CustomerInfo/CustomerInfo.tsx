import React, {useEffect, useState} from 'react';
import {Route, useParams} from 'react-router-dom';

import {Customer, FullCustomerInfo} from '@/types/customers';

import {clientFetchCustomerByID, updateCustomer, updateCustomerConfig, updateCustomerPlugins} from '@/client';

import {CenteredText} from '../CenteredText';

import {RHSTitle} from '../rhsTitle';

import {RhsHeader} from './RHSHeader';
import {CustomerInfoConfig} from './Config/Config';
import {CustomerInfoPacket} from './Packet/Packet';
import {CustomerInfoPlugins} from './Plugins/Plugins';
import {CustomerInfoProfile} from './Info/Info';
import {RhsPageHeader} from './RHSPageHeader';

interface RouteParams {
    id: string;
}

const CustomerInfo = () => {
    const {id} = useParams<RouteParams>();

    const [customer, setCustomer] = useState<FullCustomerInfo | null>(null);

    const [loading, setLoading] = useState(false);

    useEffect(() => {
        setLoading(true);
        clientFetchCustomerByID(id).
            then((res) => {
                setLoading(false);
                if (!res) {
                    return;
                }
                setCustomer(res);
            });

        return () => {
            setLoading(true);
        };
    }, [id]);

    if (loading && !customer) {
        return (
            <CenteredText
                message={'Loading...'}
            />
        );
    }

    if (!customer) {
        return (
            <CenteredText message={'Failed to pull customer info.'}/>
        );
    }

    const {config, plugins, packet, ...info} = customer;

    const handleRes = (res: FullCustomerInfo | undefined) => {
        if (!res) {
            return;
        }
        setCustomer(res);
    };

    const update = (c: Customer) => {
        updateCustomer(c.id, c).
            then(handleRes);
    };

    const saveConfig = (conf: typeof config) => {
        updateCustomerConfig(id, conf).
            then(handleRes);
    };

    const updatePlugin = (pluginID: string, isActive: boolean) => {
        updateCustomerPlugins(id, plugins.map((p) => {
            if (p.pluginID === pluginID) {
                return {
                    ...p,
                    isActive,
                };
            }
            return p;
        })).
            then(handleRes);
    };

    const deletePlugin = (pluginID: string) => {
        updateCustomerPlugins(id, plugins.filter((p) => p.pluginID !== pluginID)).
            then(handleRes);
    };

    return (
        <>
            <RHSTitle>
                <RhsHeader
                    name={info.name}
                />
            </RHSTitle>
            <RhsPageHeader
                id={info.id}
            />

            <Route path='/customers/:id/config'>
                <CustomerInfoConfig
                    config={config}
                    save={saveConfig}
                />
            </Route>
            <Route path='/customers/:id/packet'>
                <CustomerInfoPacket
                    packet={packet}
                />
            </Route>
            <Route path='/customers/:id/plugins'>
                <CustomerInfoPlugins
                    plugins={plugins}
                    changeState={updatePlugin}
                    deletePlugin={deletePlugin}
                />
            </Route>
            <Route
                path='/customers/:id'
                exact={true}
            >
                <CustomerInfoProfile
                    customer={info}
                    updateCustomer={update}
                />
            </Route>
        </>
    );
};

export {
    CustomerInfo,
};
