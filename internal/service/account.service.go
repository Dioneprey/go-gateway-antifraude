package service

import (
	"github.com/Dioneprey/go-gateway-antifraude/internal/domain"
	"github.com/Dioneprey/go-gateway-antifraude/internal/dto"
)

type AccountService struct {
	repository domain.AccountRepository
}

func NewAccountService(repository domain.AccountRepository) *AccountService {
	return &AccountService{repository: repository}
}

// CreateAccount cria uma nova conta e valida duplicidade de API Key
// Retorna ErrDuplicatedAPIKey se a chave já existir
func (s *AccountService) CreateAccount(input *dto.CreateAccountInput) (*dto.AccountOutput, error) {
	account := dto.ToAccount(input)

	// Verifica duplicidade de API Key antes da criação
	existingAccount, err := s.repository.FindByAPIKey(account.APIKey)

	if err != nil && err != domain.ErrAccountNotFound {
		return nil, err
	}
	if existingAccount != nil {
		return nil, domain.ErrDuplicatedAPIKey
	}

	err = s.repository.Save(account)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)

	return output, nil
}

// UpdateBalance atualiza o saldo de uma conta de forma thread-safe
// O amount pode ser positivo (crédito)
func (s *AccountService) UpdateBalance(apiKey string, amount float64) (*dto.AccountOutput, error) {
	account, err := s.repository.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	account.AddBalance(amount)
	err = s.repository.UpdateBalance(account)
	if err != nil {
		return nil, err
	}
	output := dto.FromAccount(account)
	return output, nil
}

// FindByAPIKey busca uma conta pelo API Key
func (s *AccountService) FindByAPIKey(apiKey string) (*dto.AccountOutput, error) {
	account, err := s.repository.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}
	output := dto.FromAccount(account)
	return output, nil
}

// FindByID busca uma conta pelo ID
func (s *AccountService) FindByID(id string) (*dto.AccountOutput, error) {
	account, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	output := dto.FromAccount(account)
	return output, nil
}
