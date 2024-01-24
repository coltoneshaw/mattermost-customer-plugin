import React, {useEffect} from 'react';

import {TextInput} from '@mantine/core';

import {UseFormReturnType, useForm} from '@mantine/form';

import {Customer, LicenseType} from '@/types/customers';
import {CenteredText} from '@/components/CenteredText';

import {Group, Container} from '@/components/Group';

import {FormUserSelector} from './ProfileSelector';
import {FormDropdown} from './CustomerTypeSelector';

export type FormTextInputParams = {
    getInputProps: UseFormReturnType<Customer, (values: Customer) => Customer>['getInputProps'];
    label: string;
    placeholder?: string;
    formKey: keyof Customer;
}

const FormTextInput = ({
    getInputProps,
    label,
    placeholder,
    formKey,
}: FormTextInputParams) => {
    return (
        <TextInput
            bg={'transparent'}
            sx={{
                width: '100%',
                label: {
                    fontSize: '14px',
                    fontWeight: 600,
                },
                input: {
                    height: '34px',
                    color: 'var(--center-channel-color)',
                    background: 'var(--center-channel-bg)',
                    border: '1px solid #ccc',
                    padding: '6px 12px',
                    borderColor: 'rgba(var(--center-channel-color-rgb), 0.16)',
                    borderRadius: '4px',
                    lineHeight: '1.42857143',
                    transition: 'border-color ease-in-out .15s, box-shadow ease-in-out .15s, -webkit-box-shadow ease-in-out .15s',
                    '&:focus': {
                        borderColor: 'rgba(var(--button-bg-rgb), 1)',
                        boxShadow: '0 0 0 1px var(--button-bg)',
                    },
                },
            }}
            placeholder={placeholder || label}
            label={label}
            {...getInputProps(formKey)}
        />
    );
};

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
                    label={isDirty() ? 'Customer Name (Unsaved Changes)' : 'Customer Name'}
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
