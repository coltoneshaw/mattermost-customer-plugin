import React, {useEffect} from 'react';

import {useForm} from '@mantine/form';

import {Customer, LicenseType} from '@/types/customers';
import {CenteredText} from '@/components/CenteredText';

import {Group, Container} from '@/components/Group';

import {FormDropdown} from '../../form/FormDropdown';

import {FormTextInput} from '@/components/form/FormTextInput';

import {FormUserSelector} from './ProfileSelector';

type Params = {
    customer: Customer | null;
    updateCustomer: (values: Customer) => void;
}
const CustomerInfoProfile = ({
    customer,
    updateCustomer,
}: Params) => {
    const {setValues, getInputProps, resetDirty, isDirty, resetTouched, values} = useForm<Customer>({
        initialValues: {
            name: '',
            accountExecutive: '',
            customerChannel: '',
            customerSuccessManager: '',
            GDriveLink: '',
            id: '',
            licensedTo: '',
            salesforceId: '',
            siteURL: '',
            technicalAccountManager: '',
            licenseType: '',
            zendeskId: '',
            lastUpdated: 0,
            airGapped: false,
            airGappedReason: '',
            codeWord: '',
            companyType: '',
            productManager: '',
            region: '',
            status: '',
        },
    });

    // resetting the form is anything new comes in.
    useEffect(() => {
        if (customer) {
            setValues(customer);
            resetDirty(customer);
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [customer]);

    if (!customer) {
        return (
            <CenteredText
                message={'No customer information found.'}
            />
        );
    }

    const resetForm = () => {
        setValues(customer);
        resetDirty(customer);
        resetTouched();
    };

    const saveValues = () => {
        updateCustomer(values);
    };

    return (
        <>
            <Container>
                <FormTextInput
                    getInputProps={getInputProps}
                    label={'Customer Name'}
                    placeholder={'Customer Name'}
                    formKey={'name'}
                />
                <FormDropdown
                    getInputProps={getInputProps}
                    label={'License Type'}
                    formKey={'licenseType'}
                    data={[
                        {value: LicenseType.Enterprise, label: 'Enterprise'},
                        {value: LicenseType.Professional, label: 'Professional'},
                        {value: LicenseType.Cloud, label: 'Cloud'},
                        {value: LicenseType.Free, label: 'Free'},
                        {value: LicenseType.Trial, label: 'Trial'},
                        {value: LicenseType.Other, label: 'Other'},
                    ]}
                />

                <FormUserSelector
                    getInputProps={getInputProps}
                    label={'Account Executive'}
                    formKey={'accountExecutive'}
                />
                <FormUserSelector
                    getInputProps={getInputProps}
                    label={'Customer Success Manager'}
                    formKey={'customerSuccessManager'}
                />
                <FormUserSelector
                    getInputProps={getInputProps}
                    label={'Technical Account Manager'}
                    formKey={'technicalAccountManager'}
                />
                <FormTextInput
                    getInputProps={getInputProps}
                    label={'Google Drive Link'}
                    formKey={'GDriveLink'}
                />
                <FormTextInput
                    getInputProps={getInputProps}
                    label={'Customer PS Channel'}
                    formKey={'customerChannel'}
                />
                <FormTextInput
                    getInputProps={getInputProps}
                    label={'Salesforce ID'}
                    formKey={'salesforceId'}
                />
                <FormTextInput
                    getInputProps={getInputProps}
                    label={'Zendesk ID'}
                    formKey={'zendeskId'}
                />
            </Container>
            {isDirty() && (
                <Group
                    style={{
                        borderTop: '1px solid rgba(var(--center-channel-color-rgb), 0.08)',
                        padding: '1em',
                    }}
                >
                    <button
                        className='btn btn-primary'
                        onClick={saveValues}
                    >
                        {'Save'}
                    </button>
                    <button
                        className='btn btn-tertiary'
                        onClick={resetForm}
                    >
                        {'Cancel'}
                    </button>
                </Group>
            )
            }

        </>

    );
};
export {
    CustomerInfoProfile,
};
