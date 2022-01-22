'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    /*
    * Helper method for generating a simplified date string.
    */
    async randomDate() {
        let year = (2020 + Math.floor(Math.random() * 10)).toString()
        let month = Math.floor((Math.random() * 12) + 1)
        let day = Math.floor((Math.random() * 31) + 1)
        if (month < 10) month = '0' + month;
        if (day < 10) day = '0' + day;
        return [year, month, day].join('.');
    }

    /**
    * Initialize the workload module with the given parameters.
    * @param {number} workerIndex The 0-based index of the worker instantiating the workload module.
    * @param {number} totalWorkers The total number of workers participating in the round.
    * @param {number} roundIndex The 0-based index of the currently executing round.
    * @param {Object} roundArguments The user-provided arguments for the round from the benchmark configuration file.
    * @param {ConnectorBase} sutAdapter The adapter of the underlying SUT.
    * @param {Object} sutContext The custom context object provided by the SUT adapter.
    * @async
    */
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);

        // Initiliase the ledger with mock data
        for (let i = 0; i < this.roundArguments.assets; i++) {

            const medNumber = `${this.workerIndex}_${i}`;

            console.log(`Worker ${this.workerIndex}: Creating asset ${medNumber}`);
            const issue = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'Issue',
                invokerIdentity: 'bob',
                contractArguments: ['Aspirin', medNumber, 'Pain Management', '2022.02.22', '$10', 'bob', 'tpmkey'],
                readOnly: false
            };
            await this.sutAdapter.sendRequests(issue);
        }
    
    }

    async submitTransaction() {
        const randomId = Math.floor(Math.random() * this.roundArguments.assets);
        const myArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'ChangeHolder',
            invokerIdentity: 'bob',
            contractArguments: ['Aspirin', `${this.workerIndex}_${randomId}`, 'charlie', 'bob', 'tpmkey'],
            readOnly: true
        };
        await this.sutAdapter.sendRequests(myArgs);
    }

    async cleanupWorkloadModule() {
        for (let i = 0; i < this.roundArguments.assets; i++) {
            const medNumber = `${this.workerIndex}_${i}`;
            console.log(`Worker ${this.workerIndex}: Deleting asset ${medNumber}`);
            const request = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'Delete',
                invokerIdentity: 'bob',
                contractArguments: ['Aspirin', medNumber, 'bob', 'tpmkey'],
                readOnly: false
            };
            await this.sutAdapter.sendRequests(request);
        }
    }

}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;