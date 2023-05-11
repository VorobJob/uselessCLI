package brain

func (b *brain) CreateTable() error {
	if err := b.repo.CreateTable(); err != nil {
		return err
	}
	return nil
}
