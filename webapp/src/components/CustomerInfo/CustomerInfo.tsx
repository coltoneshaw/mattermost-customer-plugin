import React, {useEffect, useState} from 'react';
import {Route, useParams} from 'react-router-dom';
import styled from 'styled-components';

import {FullCustomerInfo} from '@/types/customers';

import {clientFetchCustomerByID} from '@/client';

import {CenteredText} from '../CenteredText';

import {RHSTitle} from '../rhsTitle';

import {RhsHeader} from './RHSHeader';
import {CustomerInfoConfig} from './Config/Config';
import {CustomerInfoPacket} from './Packet/Packet';
import {CustomerInfoPlugins} from './Plugins/Plugins';
import {CustomerInfoProfile} from './Info/Info';
import {RhsPageHeader} from './RHSPageHeader';

const Container = styled.div`
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: flex-start;
    align-items: center;
    flex-direction: column;
    overflow-y: auto;
    gap: 1em;
`;

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

    let AlternateContainer = <CenteredText message={'Failed to pull customer info.'}/>;
    if (loading) {
        AlternateContainer = (
            <CenteredText
                message={'Loading...'}
            />
        );
    }

    if (!customer) {
        return AlternateContainer;
    }

    const {config, plugins, packet, ...info} = customer;

    return (
        <>
            <RHSTitle>
                <RhsHeader
                    name={info.name}
                />
            </RHSTitle>
            <Container>
                <RhsPageHeader
                    id={info.id}
                />
                <Route path='/customers/:id/config'>
                    <CustomerInfoConfig
                        config={config}
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
                    />
                </Route>
                <Route
                    path='/customers/:id'
                    exact={true}
                >
                    <CustomerInfoProfile
                        customer={info}
                    />
                </Route>
            </Container>
        </>
    );
};

export {
    CustomerInfo,
};
