package api

type AsaasStatus int64

const (
	AsaasPending                    AsaasStatus = iota + 1 //- Aguardando pagamento
	AsaasReceived                                          //- Recebida (saldo já creditado na conta)
	AsaasConfirmed                                         //- Pagamento confirmado (saldo ainda não creditado)
	AsaasOverdue                                           //- Vencida
	AsaasRefunded                                          //- Estornada
	AsaasReceivedInCash                                    //- Recebida em dinheiro (não gera saldo na conta)
	AsaasRefundRequested                                   //- Estorno Solicitado
	AsaasChargebackRequested                               //- Recebido chargeback
	AsaasChargebackDispute                                 //- Em disputa de chargeback (caso sejam apresentados documentos para contestação)
	AsaasAwaitingChargebackReversal                        //- Disputa vencida, aguardando repasse da adquirente
	AsaasDunningRequested                                  //- Em processo de recuperação
	AsaasDunningReceived                                   //- Recuperada
	AsaasAwaitingRiskAnalysis                              //- Pagamento em análise
	AsaasActive                                            // subscription
	AsaasExpired                                           // subscription
	AsaasDeleted
	AsaasSuccess
)
